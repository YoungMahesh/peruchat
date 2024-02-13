import Head from "next/head";

export default function PageTitle({pageName, pageDescription}: {pageName: string, pageDescription: string}) {
  return (
    <Head>
      <title>{`${pageName} | Peru Chat`}</title>
      <meta name="description" content={pageDescription} />
      <link rel="icon" href="/favicon.ico" />
    </Head>
  );
}
