import { useDeleteNotificationMutation } from "~/gen/graphql";

export const useDeleteNotification = (refetch: () => void) => {
  const [mutate] = useDeleteNotificationMutation();
  const handleDeleteNotification = () => {
    mutate().then(() => {
      refetch();
    });
  };
  return {
    handleDeleteNotification,
  };
};
