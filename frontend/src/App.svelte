<script>
  import Router, { replace } from "svelte-spa-router";
  import Home from "./pages/home.svelte";
  import SignUp from "./pages/signup.svelte";
  import NotFound from "./pages/notfound.svelte";

  const routes = {
    // Exact path
    "/": Home,
    "/signup": SignUp,
    "*": NotFound,
  };

  let sessionId = localStorage.getItem("sessionId")
  if (sessionId == null) {
    replace('/signup')
  }
</script>

<Router {routes} />
