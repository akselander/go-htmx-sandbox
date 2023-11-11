module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: ["eslint:recommended", "prettier"],
  overrides: [
    {
      env: {
        node: true,
      },
      files: [".eslintrc.{js,cjs}"],
      parserOptions: {
        sourceType: "script",
      },
    },
  ],
  plugins: ["prettier"],
  rules: {
    "prettier/prettier": "error",
  },
  ignorePatterns: ["htmx/*.js"],
};
