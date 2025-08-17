<script>
    import { apiRequest } from "../lib/api";
    export let navigate;

    async function handleLogout() {
        // Clear user id from local storage
        localStorage.removeItem("userId");

        // Make API request to logout
        try {
            const response = await apiRequest("/logout", "POST");
        } catch (error) {
            console.error("Logout failed:", error);
            // Optionally, handle the error (e.g., show a notification)
        }

        // Redirect to home page
        navigate("home");
    }
</script>

<button class="banner-btn" on:click={handleLogout}>Logout</button>
