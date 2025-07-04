import {
  GetNotificationDocument,
  useDeleteNotificationMutation,
  type GetNotificationQuery,
} from "~/gen/graphql";

export const useDeleteNotification = () => {
  const [mutate] = useDeleteNotificationMutation({
    update: (cache, { data }) => {
      const deletedId = data?.deleteNotification;
      const exist = cache.readQuery<GetNotificationQuery>({
        query: GetNotificationDocument,
      });
      if (!exist?.notification || !deletedId) {
        return;
      }
      if (exist.notification.id !== deletedId) {
        return;
      }
      cache.writeQuery({
        query: GetNotificationDocument,
        data: {
          notification: null,
        },
      });
    },
  });
  const handleDeleteNotification = () => {
    mutate();
  };
  return {
    handleDeleteNotification,
  };
};
