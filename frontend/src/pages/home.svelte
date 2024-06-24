<script>
    import { getConfig } from "../lib/api";

    let config;
    async function refreshConfig() {
        config = await getConfig()
        alert(config)
    }
</script>

<h1>Hello, World!</h1>
<button on:click={refreshConfig}>Get Config</button>