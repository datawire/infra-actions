import globals from "globals";
import js from "@eslint/js";

export default [
  {
    ...js.configs.recommended,
    languageOptions: {
      sourceType: "module",
      globals: {
        ...globals.node,
        ...globals.jest,
      },
    },
  },
];
