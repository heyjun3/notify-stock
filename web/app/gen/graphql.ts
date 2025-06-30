import { gql } from "@apollo/client";
import * as Apollo from "@apollo/client";
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = {
  [_ in K]?: never;
};
export type Incremental<T> =
  | T
  | { [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  Time: { input: string; output: string };
};

export type ChartInput = {
  end: Scalars["Time"]["input"];
  start: Scalars["Time"]["input"];
  symbol?: InputMaybe<Scalars["ID"]["input"]>;
};

export type Mutation = {
  __typename?: "Mutation";
  createNotification: Notification;
  deleteNotification: Scalars["ID"]["output"];
};

export type MutationCreateNotificationArgs = {
  input: NotificationInput;
};

export type Node = {
  id: Scalars["ID"]["output"];
};

export type Notification = Node & {
  __typename?: "Notification";
  id: Scalars["ID"]["output"];
  targets: Array<SymbolDetail>;
  time: Scalars["Time"]["output"];
};

export type NotificationInput = {
  symbols: Array<Scalars["ID"]["input"]>;
  time: Scalars["Time"]["input"];
};

export type Query = {
  __typename?: "Query";
  node?: Maybe<Node>;
  notification?: Maybe<Notification>;
  notifications: Array<Notification>;
  symbol: Symbol;
  symbols: Array<Symbol>;
};

export type QueryNodeArgs = {
  id: Scalars["ID"]["input"];
};

export type QuerySymbolArgs = {
  input: SymbolInput;
};

export type QuerySymbolsArgs = {
  input?: InputMaybe<SymbolInput>;
};

export type Stock = {
  __typename?: "Stock";
  price: Scalars["Float"]["output"];
  symbol: Scalars["ID"]["output"];
  timestamp: Scalars["String"]["output"];
};

export type Symbol = Node & {
  __typename?: "Symbol";
  chart: Array<Stock>;
  detail: SymbolDetail;
  id: Scalars["ID"]["output"];
  symbol: Scalars["ID"]["output"];
};

export type SymbolChartArgs = {
  input: ChartInput;
};

export type SymbolDetail = {
  __typename?: "SymbolDetail";
  change: Scalars["String"]["output"];
  changePercent: Scalars["String"]["output"];
  currencySymbol: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  longName: Scalars["String"]["output"];
  marketCap?: Maybe<Scalars["String"]["output"]>;
  price: Scalars["Float"]["output"];
  shortName: Scalars["String"]["output"];
  symbol: Scalars["ID"]["output"];
  volume?: Maybe<Scalars["String"]["output"]>;
};

export type SymbolInput = {
  symbol: Scalars["ID"]["input"];
};

export type CreateNotificationMutationVariables = Exact<{
  createNotificationInput: NotificationInput;
}>;

export type CreateNotificationMutation = {
  __typename?: "Mutation";
  createNotification: {
    __typename?: "Notification";
    id: string;
    time: string;
    targets: Array<{ __typename?: "SymbolDetail"; id: string; symbol: string; shortName: string }>;
  };
};

export type DeleteNotificationMutationVariables = Exact<{ [key: string]: never }>;

export type DeleteNotificationMutation = { __typename?: "Mutation"; deleteNotification: string };

export type GetSymbolsQueryVariables = Exact<{
  chartInput: ChartInput;
}>;

export type GetSymbolsQuery = {
  __typename?: "Query";
  symbols: Array<{
    __typename?: "Symbol";
    id: string;
    symbol: string;
    detail: {
      __typename?: "SymbolDetail";
      symbol: string;
      shortName: string;
      longName: string;
      price: number;
      change: string;
      changePercent: string;
      volume?: string | null;
      marketCap?: string | null;
      currencySymbol: string;
    };
    chart: Array<{ __typename?: "Stock"; symbol: string; timestamp: string; price: number }>;
  }>;
};

export type GetNotificationQueryVariables = Exact<{ [key: string]: never }>;

export type GetNotificationQuery = {
  __typename?: "Query";
  notification?: {
    __typename?: "Notification";
    id: string;
    time: string;
    targets: Array<{ __typename?: "SymbolDetail"; id: string; symbol: string; shortName: string }>;
  } | null;
};

export const CreateNotificationDocument = gql`
    mutation createNotification($createNotificationInput: NotificationInput!) {
  createNotification(input: $createNotificationInput) {
    id
    time
    targets {
      id
      symbol
      shortName
    }
  }
}
    `;
export type CreateNotificationMutationFn = Apollo.MutationFunction<
  CreateNotificationMutation,
  CreateNotificationMutationVariables
>;

/**
 * __useCreateNotificationMutation__
 *
 * To run a mutation, you first call `useCreateNotificationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateNotificationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createNotificationMutation, { data, loading, error }] = useCreateNotificationMutation({
 *   variables: {
 *      createNotificationInput: // value for 'createNotificationInput'
 *   },
 * });
 */
export function useCreateNotificationMutation(
  baseOptions?: Apollo.MutationHookOptions<
    CreateNotificationMutation,
    CreateNotificationMutationVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useMutation<CreateNotificationMutation, CreateNotificationMutationVariables>(
    CreateNotificationDocument,
    options,
  );
}
export type CreateNotificationMutationHookResult = ReturnType<typeof useCreateNotificationMutation>;
export type CreateNotificationMutationResult = Apollo.MutationResult<CreateNotificationMutation>;
export type CreateNotificationMutationOptions = Apollo.BaseMutationOptions<
  CreateNotificationMutation,
  CreateNotificationMutationVariables
>;
export const DeleteNotificationDocument = gql`
    mutation deleteNotification {
  deleteNotification
}
    `;
export type DeleteNotificationMutationFn = Apollo.MutationFunction<
  DeleteNotificationMutation,
  DeleteNotificationMutationVariables
>;

/**
 * __useDeleteNotificationMutation__
 *
 * To run a mutation, you first call `useDeleteNotificationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteNotificationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteNotificationMutation, { data, loading, error }] = useDeleteNotificationMutation({
 *   variables: {
 *   },
 * });
 */
export function useDeleteNotificationMutation(
  baseOptions?: Apollo.MutationHookOptions<
    DeleteNotificationMutation,
    DeleteNotificationMutationVariables
  >,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useMutation<DeleteNotificationMutation, DeleteNotificationMutationVariables>(
    DeleteNotificationDocument,
    options,
  );
}
export type DeleteNotificationMutationHookResult = ReturnType<typeof useDeleteNotificationMutation>;
export type DeleteNotificationMutationResult = Apollo.MutationResult<DeleteNotificationMutation>;
export type DeleteNotificationMutationOptions = Apollo.BaseMutationOptions<
  DeleteNotificationMutation,
  DeleteNotificationMutationVariables
>;
export const GetSymbolsDocument = gql`
    query GetSymbols($chartInput: ChartInput!) {
  symbols {
    id
    symbol
    detail {
      symbol
      shortName
      longName
      price
      change
      changePercent
      volume
      marketCap
      currencySymbol
    }
    chart(input: $chartInput) {
      symbol
      timestamp
      price
    }
  }
}
    `;

/**
 * __useGetSymbolsQuery__
 *
 * To run a query within a React component, call `useGetSymbolsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSymbolsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSymbolsQuery({
 *   variables: {
 *      chartInput: // value for 'chartInput'
 *   },
 * });
 */
export function useGetSymbolsQuery(
  baseOptions: Apollo.QueryHookOptions<GetSymbolsQuery, GetSymbolsQueryVariables> &
    ({ variables: GetSymbolsQueryVariables; skip?: boolean } | { skip: boolean }),
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<GetSymbolsQuery, GetSymbolsQueryVariables>(GetSymbolsDocument, options);
}
export function useGetSymbolsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<GetSymbolsQuery, GetSymbolsQueryVariables>,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<GetSymbolsQuery, GetSymbolsQueryVariables>(
    GetSymbolsDocument,
    options,
  );
}
export function useGetSymbolsSuspenseQuery(
  baseOptions?:
    | Apollo.SkipToken
    | Apollo.SuspenseQueryHookOptions<GetSymbolsQuery, GetSymbolsQueryVariables>,
) {
  const options =
    baseOptions === Apollo.skipToken ? baseOptions : { ...defaultOptions, ...baseOptions };
  return Apollo.useSuspenseQuery<GetSymbolsQuery, GetSymbolsQueryVariables>(
    GetSymbolsDocument,
    options,
  );
}
export type GetSymbolsQueryHookResult = ReturnType<typeof useGetSymbolsQuery>;
export type GetSymbolsLazyQueryHookResult = ReturnType<typeof useGetSymbolsLazyQuery>;
export type GetSymbolsSuspenseQueryHookResult = ReturnType<typeof useGetSymbolsSuspenseQuery>;
export type GetSymbolsQueryResult = Apollo.QueryResult<GetSymbolsQuery, GetSymbolsQueryVariables>;
export const GetNotificationDocument = gql`
    query getNotification {
  notification {
    id
    time
    targets {
      id
      symbol
      shortName
    }
  }
}
    `;

/**
 * __useGetNotificationQuery__
 *
 * To run a query within a React component, call `useGetNotificationQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNotificationQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNotificationQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNotificationQuery(
  baseOptions?: Apollo.QueryHookOptions<GetNotificationQuery, GetNotificationQueryVariables>,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<GetNotificationQuery, GetNotificationQueryVariables>(
    GetNotificationDocument,
    options,
  );
}
export function useGetNotificationLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<GetNotificationQuery, GetNotificationQueryVariables>,
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<GetNotificationQuery, GetNotificationQueryVariables>(
    GetNotificationDocument,
    options,
  );
}
export function useGetNotificationSuspenseQuery(
  baseOptions?:
    | Apollo.SkipToken
    | Apollo.SuspenseQueryHookOptions<GetNotificationQuery, GetNotificationQueryVariables>,
) {
  const options =
    baseOptions === Apollo.skipToken ? baseOptions : { ...defaultOptions, ...baseOptions };
  return Apollo.useSuspenseQuery<GetNotificationQuery, GetNotificationQueryVariables>(
    GetNotificationDocument,
    options,
  );
}
export type GetNotificationQueryHookResult = ReturnType<typeof useGetNotificationQuery>;
export type GetNotificationLazyQueryHookResult = ReturnType<typeof useGetNotificationLazyQuery>;
export type GetNotificationSuspenseQueryHookResult = ReturnType<
  typeof useGetNotificationSuspenseQuery
>;
export type GetNotificationQueryResult = Apollo.QueryResult<
  GetNotificationQuery,
  GetNotificationQueryVariables
>;
