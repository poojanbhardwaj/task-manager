import { Routes, Route } from "react-router-dom";
import { SignIn, SignUp } from "@clerk/react";
import { dark } from "@clerk/themes";
import { useTheme } from "@/components/theme-provider";

export function AuthLayout() {
  const { theme } = useTheme();
  const isDark = theme === "dark" || (theme === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches);

  return (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <div className="w-full max-w-md">
        <Routes>
          <Route
            path="/sign-in/*"
            element={
              <SignIn
                fallbackRedirectUrl="/dashboard"
                appearance={{
                  theme: isDark ? dark : undefined,
                  elements: {
                    card: "shadow-lg border-border",
                    rootBox: "w-full",
                  },
                }}
              />
            }
          />
          <Route
            path="/sign-up/*"
            element={
              <SignUp
                fallbackRedirectUrl="/dashboard"
                appearance={{
                  theme: isDark ? dark : undefined,
                  elements: {
                    card: "shadow-lg border-border",
                    rootBox: "w-full",
                  },
                }}
              />
            }
          />
        </Routes>
      </div>
    </div>
  );
}