<script>
    import { replace } from "svelte-spa-router";

    const params = new URLSearchParams(location.search)
    if (params.get("sessionId") !== null) {
        localStorage.setItem("sessionId", params.get("sessionId"))
        replace("/")
    }
</script>