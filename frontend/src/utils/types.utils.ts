export type FetchedMessage = {
  is_sender: boolean;
  message: string;
};

export type LoggedInUserInfo = {
  username: string;
  token: string;
}