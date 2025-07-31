<script>
    export let navigate;

    function handleLogout() {
        // Clear user id from local storage
        localStorage.removeItem("userId");

        // Redirect to home page
        navigate("home");
    }
</script>

<button on:click={handleLogout}>Logout</button>
