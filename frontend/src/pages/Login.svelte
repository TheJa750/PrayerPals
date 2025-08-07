<script>
    import { apiRequest } from "../lib/api";
    export let navigate;

    let errorMessage = "";
    let successMessage = "";
    let email = "";
    let password = "";
    let isSubmitting = false;
    let showPassword = false;

    async function login(event) {
        event.preventDefault();
        isSubmitting = true;
        errorMessage = "";
        successMessage = "";

        try {
            const loginData = {
                email: email,
                password: password,
            };

            const response = await apiRequest("/login", "POST", loginData);

            // Successfully logged in, navigate to home
            localStorage.setItem("userId", response.id);
            successMessage = "Logged in! Redirecting to your dashboard...";

            // Clear form fields
            email = "";
            password = "";

            // Navigate to user dashboard after brief delay
            setTimeout(() => {
                navigate("user");
            }, 1000);
        } catch (error) {
            errorMessage =
                error.message || "Login failed. Please check your credentials.";
            console.error("Login error:", error);
        } finally {
            isSubmitting = false;
        }
    }
</script>

<h2 class="text-center mb-1">Log In</h2>
<form on:submit={login}>
    <div class="form-row">
        <label for="login-email">Email:</label>
        <input
            type="email"
            id="login-email"
            required
            bind:value={email}
            disabled={isSubmitting}
        />
    </div>
    <div class="form-row">
        <label for="login-password">Password:</label>
        <input
            type={showPassword ? "text" : "password"}
            id="login-password"
            required
            bind:value={password}
            disabled={isSubmitting}
        />
    </div>

    <!-- Messages will go here between password and buttons -->
    {#if errorMessage || successMessage}
        <div class="form-row">
            {#if errorMessage}
                <div class="error-message text-center">{errorMessage}</div>
            {/if}
            {#if successMessage}
                <div class="success-message text-center">{successMessage}</div>
            {/if}
        </div>
    {/if}

    <div class="form-row">
        <button
            type="button"
            class="toggle-password"
            on:click={() => (showPassword = !showPassword)}
            aria-label={showPassword ? "Hide password" : "Show password"}
            disabled={isSubmitting}
        >
            {showPassword ? "🙈 Hide" : "👁️ Show"}
        </button>
        <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? "Logging in..." : "Log In"}
        </button>
    </div>
</form>
