import { ApolloError } from "@apollo/client";
import { parse, format } from "@formkit/tempo";

import { useGetNotificationQuery } from "~/gen/graphql";

const toLocalTime = (time: string): string => {
  const date = parse(time);
  return format(date, "HH:mm", "jp");
};

const isError = (error?: ApolloError) => {
  if (error === undefined) return false;
  const extension = error.cause?.extensions;
  const extensions = Array.isArray(extension) ? extension : extension == null ? [] : [extension];
  return extensions.find((e) => e.code) ? true : false;
};
export const useGetNotification = () => {
  const { data, loading, error } = useGetNotificationQuery();
  let notifications = undefined;
  if (data?.notification) {
    notifications = [
      {
        id: data.notification.id,
        time: toLocalTime(data.notification.time),
        tickers: data.notification.targets.map((t) => t.shortName),
      },
    ];
  }
  return {
    notifications: notifications ?? [],
    loading,
    unAuthorization: isError(error),
  };
};
