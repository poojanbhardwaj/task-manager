import "./index.css";
import { Show, SignInButton, SignUpButton, UserButton } from "@clerk/react";

function App() {
  return (
    <div className="mx-auto max-w-4xl p-6">
      <header className="flex flex-wrap items-center justify-between gap-4">
        <h1 className="text-2xl font-semibold">Task Manager</h1>
        <nav aria-label="Account" className="flex items-center gap-3">
          <Show when="signed-out">
            <SignInButton mode="modal"><button type="button" className="rounded-md border px-4 py-2 text-sm font-medium">Sign in</button></SignInButton>
            <SignUpButton mode="modal"><button type="button" className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground">Create account</button></SignUpButton>
          </Show>
          <Show when="signed-in"><UserButton /></Show>
        </nav>
      </header>
      <main className="mt-8">
        <Show when="signed-out"><p className="text-muted-foreground">Sign in to get started.</p></Show>
        <Show when="signed-in"><p className="text-muted-foreground">You are signed in.</p></Show>
      </main>
    </div>
  );
}

export default App;