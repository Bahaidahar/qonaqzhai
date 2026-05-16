import type { Metadata } from "next";
import { Manrope } from "next/font/google";
import "./globals.css";
import { I18nProvider } from "@/shared/i18n/context";
import { SidebarProvider } from "@/widgets/chat-sidebar/sidebar-context";
import { ThemeProvider } from "@/features/theme/context";
import { AuthProvider } from "@/features/auth/context";

const THEME_INIT = `(function(){try{var t=localStorage.getItem('qonaqzhai_theme')||'light';var r=t==='system'?(window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light'):t;document.documentElement.setAttribute('data-theme',r);}catch(e){}})();`;

const manrope = Manrope({
  variable: "--font-sans-raw",
  subsets: ["latin", "cyrillic", "cyrillic-ext"],
  display: "swap",
});

export const metadata: Metadata = {
  title: "qonaqzhai — Plan events by chatting",
  description:
    "AI-powered event planning for weddings, toi, corporate, birthdays. Generate timelines, budgets, vendor matches by chatting.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${manrope.variable} h-full antialiased`}
      suppressHydrationWarning
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: THEME_INIT }} />
      </head>
      <body className="min-h-full flex flex-col">
        <ThemeProvider>
          <I18nProvider>
            <AuthProvider>
              <SidebarProvider>{children}</SidebarProvider>
            </AuthProvider>
          </I18nProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
