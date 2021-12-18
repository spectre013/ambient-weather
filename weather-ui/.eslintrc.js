module.exports = {
  root: true,
  env: {
    node: true,
    browser: true,
  },
  extends: ['plugin:vue/essential', 'plugin:prettier/recommended', '@vue/prettier'],
  rules: {
    'eqeqeq': 'off',
    'no-plusplus': 'off',
    "no-return-assign": "off",
    'max-len': ['error', { code: 120, ignoreUrls: true, ignoreComments: true }],
    'linebreak-style': 0,
    'vue/max-attributes-per-line': 'off',
    'vue/component-name-in-template-casing': ['error', 'PascalCase'],
    'no-console': 'off',
    'no-restricted-syntax': [
      'error',
      {
        selector:
            "CallExpression[callee.object.name='console'][callee.property.name!=/^(log|warn|error|info|trace)$/]",
        message: 'Unexpected property on console object was called',
      },
    ],
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
};