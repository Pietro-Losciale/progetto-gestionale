import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

// import bootstrap italia
import "bootstrap-italia/dist/css/bootstrap-italia.min.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata = {
  title: "Gestionale Magazzino",
  description: "Progetto gestionale full stack con Next.js e Go",
};

export default function RootLayout({ children }) {
  return (
    <html
      lang="it"
      className={`${geistSans.variable} ${geistMono.variable}`}
    >
      <body>
        {children}
      </body>
    </html>
  );
}
