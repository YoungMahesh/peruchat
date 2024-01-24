import { type AppType } from "next/dist/shared/lib/utils";
import { NextUIProvider } from "@nextui-org/react";
import "~/styles/globals.css";

const MyApp: AppType = ({ Component, pageProps }) => {
  return (
    <NextUIProvider>
      <Component {...pageProps} />
    </NextUIProvider>
  );
};

export default MyApp;
