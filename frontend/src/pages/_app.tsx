import { type AppType } from "next/dist/shared/lib/utils";
import { NextUIProvider } from "@nextui-org/react";
import { ToastContainer } from "react-toastify";
import "~/styles/globals.css";
import "react-toastify/dist/ReactToastify.css";

const MyApp: AppType = ({ Component, pageProps }) => {
  return (
    <NextUIProvider>
      <Component {...pageProps} />
      <ToastContainer />
    </NextUIProvider>
  );
};

export default MyApp;
