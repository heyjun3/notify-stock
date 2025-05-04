import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  overwrite: true,
  schema: "http://localhost:8080/query",
  documents: "app/**/*.gql",
  generates: {
    "app/gen/graphql.ts": {
      plugins: ["typescript", "typescript-operations", "typescript-react-apollo"],
      config: {
        withHooks: true,
        withComponent: false,
        scalars: {
          Time: "string",
        },
      },
    },
  },
};
export default config;
