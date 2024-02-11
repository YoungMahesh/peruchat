export type FetchedMessage = {
  is_sender: boolean;
  message: string;
};

export type Event1 = {
  type: string;
  payload: unknown;
}

export type LoggedInUserInfo = {
  username: string;
  token: string;
}