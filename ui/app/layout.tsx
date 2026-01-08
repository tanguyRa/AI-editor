import type { Metadata } from "next";
import { Montserrat } from "next/font/google";
import "./globals.css";

const appFont = Montserrat({
  subsets: ["latin"],
  variable: "--font-app",
  weight: ["400", "700"],
});

export const metadata: Metadata = {
  title: "AI editor",
  description: "UI generator using AI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={appFont.className}
      >
        {children}
      </body>
    </html>
  );
}
