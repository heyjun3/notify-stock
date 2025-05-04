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
};

export type MutationCreateNotificationArgs = {
  input: NotificationInput;
};

export type Notification = {
  __typename?: "Notification";
  email: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  symbol: Scalars["ID"]["output"];
  time: Scalars["Time"]["output"];
};

export type NotificationInput = {
  email: Scalars["String"]["input"];
  symbol: Scalars["ID"]["input"];
  time: Scalars["Time"]["input"];
};

export type Query = {
  __typename?: "Query";
  symbol: Symbol;
  symbols: Array<Symbol>;
};

export type QuerySymbolArgs = {
  input: SymbolInput;
};

export type QuerySymbolsArgs = {
  input?: InputMaybe<SymbolInput>;
};

export type Stock = {
  __typename?: "Stock";
  close: Scalars["Float"]["output"];
  symbol: Scalars["ID"]["output"];
  timestamp: Scalars["Time"]["output"];
};

export type Symbol = {
  __typename?: "Symbol";
  chart?: Maybe<Array<Maybe<Stock>>>;
  currentStock: Stock;
  detail: SymbolDetail;
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

export type GetSymbolsQueryVariables = Exact<{
  chartInput: ChartInput;
}>;

export type GetSymbolsQuery = {
  __typename?: "Query";
  symbols: Array<{
    __typename?: "Symbol";
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
    chart?: Array<{
      __typename?: "Stock";
      symbol: string;
      timestamp: string;
      close: number;
    } | null> | null;
  }>;
};

export const GetSymbolsDocument = gql`
    query GetSymbols($chartInput: ChartInput!) {
  symbols {
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
      close
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
