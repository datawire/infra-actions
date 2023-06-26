import os
import requests
import json
from datetime import datetime, timedelta, timezone, date

GITHUB_API_URL = "https://api.github.com"
REPO_OWNER = "telepresenceio"
REPO_NAME = "telepresence"

def dump_to_file(data, filename):
    with open(filename, "w") as f:
        json.dump(data, f)

def dict_from_json_file(filename):
    with open(filename, "r") as f:
        return json.load(f)


def fetch_workflow_runs(token, date_threshold, days):
    print(f"Fetching workflow runs from {date_threshold} to {days} days ago.")

    headers = {
        "Authorization": f"Bearer {token}",
        "Accept": "application/vnd.github+json"
    }
    url = f"{GITHUB_API_URL}/repos/{REPO_OWNER}/{REPO_NAME}/actions/runs"

    params = {
        "status": "completed",
        "created": f">{date_threshold}"
    }

    if os.path.exists(f"workflow_runs-{days}days-to-{date_threshold}.json"):
        return dict_from_json_file(f"workflow_runs-{days}days-to-{date_threshold}.json")

    workflow_runs = []
    page = 1
    while True:
        params["page"] = page
        response = requests.get(url, headers=headers, params=params)
        response.raise_for_status()
        page_workflow_runs = response.json()["workflow_runs"]
        workflow_runs.extend(page_workflow_runs)
        page += 1
        if len(page_workflow_runs) == 0:
            break

    dump_to_file(workflow_runs, f"workflow_runs-{days}days-to-{date_threshold}.json")

    return workflow_runs

def fetch_workflow_jobs(token, workflow):
    headers = {
        "Authorization": f"Bearer {token}",
        "Accept": "application/vnd.github+json"
    }
    url = workflow.get("jobs_url")
    params = {
        "per_page": 100,
        "filter": "all"
    }

    response = requests.get(url, headers=headers)
    response.raise_for_status()
    return response.json()["jobs"]

def get_dict_by_id(list_of_dicts, id):
    return next((item for item in list_of_dicts if item["id"] == id), None)

def fetch_jobs(token, date_threshold, days):
    workflows = fetch_workflow_runs(token, date_threshold, days)
    print(f"Found {len(workflows)} workflows in the last {days} days.")
    print(f"Fetching jobs for each workflow.")
    jobs = []

    if os.path.exists(f"jobs-{days}days-to-{date_threshold}.json"):
        jobs = dict_from_json_file(f"jobs-{days}days-to-{date_threshold}.json")
    else:
        for workflow in workflows:
            workflow_jobs = fetch_workflow_jobs(token, workflow)
            for job in workflow_jobs:
                job_workflow = get_dict_by_id(workflows, job["run_id"])
                job["workflow_created_at"] = job_workflow["created_at"]
            jobs.extend(workflow_jobs)

    dump_to_file(jobs, f"jobs-{days}days-to-{date_threshold}.json")

    return jobs


def filter_jobs(jobs, label):
    return [job for job in jobs if label in job.get("labels")]


def get_job_times(job):
    times = {}
    simple_datetime_format_string = "%Y-%m-%dT%H:%M:%SZ"
    times["created_at"] = datetime.strptime(job["created_at"], simple_datetime_format_string)
    times["started_at"] = datetime.strptime(job["started_at"], simple_datetime_format_string)
    times["completed_at"] = datetime.strptime(job["completed_at"], simple_datetime_format_string)
    times["workflow_created_at"] = datetime.strptime(job["workflow_created_at"], simple_datetime_format_string)

    fractional_datetime_format_string = "%Y-%m-%dT%H:%M:%S.%f%z"
    times["first_step_started_at"] = datetime.strptime(job["steps"][0]["started_at"], fractional_datetime_format_string).astimezone(timezone.utc)
    times["last_step_completed_at"] = datetime.strptime(job["steps"][-1]["completed_at"], fractional_datetime_format_string).astimezone(timezone.utc)

    return times


def calculate_total_minutes(jobs):
    total_minutes = 0

    for job in jobs:
        job_times = get_job_times(job)
        duration = (job_times["completed_at"] -
                    job_times["started_at"]).total_seconds() / 60
        total_minutes += duration

    return total_minutes

def calculate_total_delay(jobs):
    total_delay = 0

    for job in jobs:
        job_times = get_job_times(job)
        delay = (job_times["started_at"] -
                    job_times["workflow_created_at"]).total_seconds() / 60
        total_delay += delay

    return total_delay

def find_worst_delay(jobs):
    worst_delay = 0
    worst_job = None

    for job in jobs:
        job_times = get_job_times(job)
        delay = (job_times["started_at"] -
                    job_times["workflow_created_at"]).total_seconds() / 60
        if delay > worst_delay:
            worst_delay = delay
            worst_job = job

    return worst_delay, worst_job

def get_total_label_minutes(label, days):
    gh_token = os.environ.get("GITHUB_TOKEN")
    if not gh_token:
        raise ValueError("Environment variable 'GH_TOKEN' not set.")

    date_threshold = (date(2023, 6, 23) - timedelta(days=days+1)).isoformat()

    jobs = fetch_jobs(gh_token, date_threshold, days)

    label = "macOS-arm64"
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
    get_total_label_minutes("macOS-arm64", 30)
