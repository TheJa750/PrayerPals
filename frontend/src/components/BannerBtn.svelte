<script>
    export let navigate;
    export let page;
    export let text;
</script>

<button class="banner-btn" on:click={() => navigate(page)}>{text}</button>
