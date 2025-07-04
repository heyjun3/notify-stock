import { GetNotificationDocument, useCreateNotificationMutation } from "~/gen/graphql";

const parseTime = (time: string) => {
  const hour = Number(time.slice(0, 2));
  const date = new Date();
  date.setHours(hour);
  return date;
};

export const useCreateNotification = () => {
  const [mutate] = useCreateNotificationMutation({
    update: (cache, { data }) => {
      const newNotification = data?.createNotification;
      if (!newNotification) {
        return;
      }
      cache.writeQuery({
        query: GetNotificationDocument,
        data: {
          notification: newNotification,
        },
      });
    },
  });
  const handleCreateNotification = (notification: {
    id: number;
    time: string;
    selectedStockSymbols: string[];
  }) => {
    mutate({
      variables: {
        createNotificationInput: {
          symbols: notification.selectedStockSymbols,
          time: parseTime(notification.time).toISOString(),
        },
      },
    });
  };
  return { handleCreateNotification };
};
