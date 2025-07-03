import { ApolloError } from "@apollo/client";
import { useGetNotificationQuery } from "~/gen/graphql";

const isError = (error?: ApolloError) => {
  if (error === undefined) return false;
  const extension = error.cause?.extensions;
  const extensions = Array.isArray(extension) ? extension : extension == null ? [] : [extension];
  return extensions.find((e) => e.code) ? true : false;
};
export const useGetNotification = () => {
  const { data, loading, error, refetch } = useGetNotificationQuery();
  console.warn(data?.notification);
  let notifications = undefined;
  if (data?.notification) {
    notifications = [
      {
        id: data.notification.id,
        time: data.notification.time,
        tickers: data.notification.targets.map((t) => t.shortName),
      },
    ];
  }
  return {
    refetch,
    notifications: notifications ?? [],
    loading,
    unAuthorization: isError(error),
  };
};
