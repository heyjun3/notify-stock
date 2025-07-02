import { useCreateNotificationMutation } from "~/gen/graphql";

export const useCreateNotification = (refetch: () => void) => {
  const parseTime = (time: string) => {
    const hour = Number(time.slice(0, 2));
    const date = new Date();
    date.setHours(hour);
    return date;
  };
  const [mutate] = useCreateNotificationMutation();
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
    }).then((r) => {
      console.warn(r.data);
      refetch();
    });
  };
  return { handleCreateNotification };
};