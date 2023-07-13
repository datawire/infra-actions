import os
import requests
import json
import argparse
from datetime import datetime, timedelta, timezone, date

GITHUB_API_URL = "https://api.github.com"
REPO_OWNER = ""
REPO_NAME = ""
github_token = os.environ.get("GITHUB_TOKEN")
created_at = "created_at"
started_at = "started_at"
completed_at = "completed_at"
workflow_created_at = "workflow_created_at"

date_format_string = "%Y-%m-%d"
simple_datetime_format_string = "%Y-%m-%dT%H:%M:%SZ"
fractional_datetime_format_string = "%Y-%m-%dT%H:%M:%S.%f%z"

github_request_headers = {
    "Authorization": f"Bearer {github_token}",
    "Accept": "application/vnd.github+json"
}

def dump_to_file(data, filename):
    with open(filename, "w") as f:
        json.dump(data, f)

def dict_from_json_file(filename):
    with open(filename, "r") as f:
        return json.load(f)


def fetch_workflow_runs(date_threshold, days):
    print(f"Fetching workflow runs from {date_threshold} to {days} days ago.")

    json_filepath = f"{REPO_OWNER}_{REPO_NAME}-workflow_runs-{days}-days-to-{date_threshold.split('T')[0]}.json"
    if os.path.exists(json_filepath):
        return dict_from_json_file(json_filepath)

    url = f"{GITHUB_API_URL}/repos/{REPO_OWNER}/{REPO_NAME}/actions/runs"

    params = {
        "status": "completed",
        "created": f">{date_threshold}"
    }

    workflow_runs = []
    page = 1
    while True:
        params["page"] = page
        response = requests.get(url, headers=github_request_headers, params=params)
        response.raise_for_status()
        page_workflow_runs = response.json()["workflow_runs"]
        workflow_runs.extend(page_workflow_runs)
        page += 1
        if len(page_workflow_runs) == 0:
            break

    dump_to_file(workflow_runs, json_filepath)

    return workflow_runs

def fetch_workflow_jobs(workflow):
    url = workflow.get("jobs_url")
    params = {
        "per_page": 100,
        "filter": "all"
    }

    response = requests.get(url, headers=github_request_headers)
    response.raise_for_status()
    return response.json()["jobs"]

def get_dict_by_id(list_of_dicts, id):
    return next((item for item in list_of_dicts if item["id"] == id), None)

def fetch_jobs(label, date_threshold, days):
    workflows = fetch_workflow_runs(date_threshold, days)
    print(f"Found {len(workflows)} workflows in the last {days} days.")
    print(f"Fetching jobs for each workflow.")
    jobs = []

    json_filepath = f"{REPO_OWNER}_{REPO_NAME}-{label}-jobs-{days}days-to-{date_threshold.split('T')[0]}.json"
    if os.path.exists(json_filepath):
        jobs = dict_from_json_file(json_filepath)
    else:
        for workflow in workflows:
            workflow_jobs = fetch_workflow_jobs(workflow)
            for job in workflow_jobs:
                job_workflow = get_dict_by_id(workflows, job["run_id"])
                job["workflow_created_at"] = job_workflow["created_at"]
            jobs.extend(workflow_jobs)

    dump_to_file(jobs, json_filepath)

    return jobs


def filter_jobs(jobs, label):
    return [job for job in jobs if label in job.get("labels")]


def get_job_times(job):
    times = {}

    times[created_at] = datetime.strptime(job[created_at], simple_datetime_format_string)
    times[started_at] = datetime.strptime(job[started_at], simple_datetime_format_string)
    times[completed_at] = datetime.strptime(job[completed_at], simple_datetime_format_string)
    times[workflow_created_at] = datetime.strptime(job[workflow_created_at], simple_datetime_format_string)

    if len(job["steps"]) == 0:
        times["first_step_started_at"] = 0
        times["last_step_completed_at"] = 0
    else:
        times["first_step_started_at"] = datetime.strptime(job["steps"][0][started_at], fractional_datetime_format_string).astimezone(timezone.utc)
        times["last_step_completed_at"] = datetime.strptime(job["steps"][-1][completed_at], fractional_datetime_format_string).astimezone(timezone.utc)

    return times


def calculate_total_minutes(jobs):
    total_minutes = 0

    for job in jobs:
        job_times = get_job_times(job)
        duration = (job_times[completed_at] -
                    job_times[started_at]).total_seconds() / 60
        total_minutes += duration

    return total_minutes

def calculate_total_delay(jobs):
    total_delay = 0

    for job in jobs:
        job_times = get_job_times(job)
        delay = (job_times[started_at] -
                    job_times[workflow_created_at]).total_seconds() / 60
        total_delay += delay

    return total_delay

def find_worst_delay(jobs):
    worst_delay = 0
    worst_job = None

    for job in jobs:
        job_times = get_job_times(job)
        delay = (job_times[started_at] -
                    job_times[workflow_created_at]).total_seconds() / 60
        if delay > worst_delay:
            worst_delay = delay
            worst_job = job

    return worst_delay, worst_job

def get_total_label_minutes(labels, target_date, days):
    date_threshold = (target_date - timedelta(days=days+1)).isoformat()

    for label in labels:
        print(f"Fetching {label} jobs from {date_threshold} to {days} days before.")

        jobs = fetch_jobs(label, date_threshold, days)
        filtered_jobs = filter_jobs(jobs, label)
        num_jobs = len(filtered_jobs)
        print(f"Found {num_jobs} {label} jobs in the last {days} days.")

        total_minutes = calculate_total_minutes(filtered_jobs)
        total_delay = calculate_total_delay(filtered_jobs)
        worst_delay = find_worst_delay(filtered_jobs)
        print(f"Total {label} runtime in the last {days} days: {total_minutes:.2f} minutes.")
        print(f"Total {label} delay in the last {days} days: {total_delay} minutes.")
        print(f"Average {label} delay in the last {days} days: {total_delay/num_jobs:.2f} minutes.")
        print(f"Worst {label} delay in the last {days} days: {worst_delay[0]:.2f} minutes.")



if __name__ == "__main__":
    if not github_token:
        raise ValueError("Environment variable 'GH_TOKEN' not set.")

    parser = argparse.ArgumentParser(description="Fetch and analyze GitHub Actions job data")
    parser.add_argument("command", choices=["total_minutes"], help="Command to run")
    parser.add_argument("repo", type=str, help="Repository to fetch data for in the format 'owner/repo'")
    parser.add_argument("--target-date", type=str, default=date.today().strftime(date_format_string), help="Most recent date to fetch data for. YYYY-MM-DD format")
    parser.add_argument("--num-days", type=int, default=30, help="Number of days before the target date to fetch data for")
    parser.add_argument("--labels", type=str, nargs="+", required=True, help="Labels to filter jobs by")

    args = parser.parse_args()

    REPO_OWNER = args.repo.split("/")[0]
    REPO_NAME = args.repo.split("/")[1]

    target_date = datetime.strptime(args.target_date, date_format_string)
    num_days = args.num_days
    labels = args.labels

    if args.command == "total_minutes":
        get_total_label_minutes(labels, target_date, args.num_days)
