// src/lib/api.js
export const BASE_URL = "http://localhost:8080/api"; // Adjust when deploying
export const REFRESH_ERROR_MESSAGE = "Session expired. Please log in again";

export async function apiRequest(endpoint, method = "GET", data = undefined) {
    const url = BASE_URL + endpoint;
    const options = {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include" // needed for cookies/auth, adjust as needed
    };
    if (data) {
        options.body = JSON.stringify(data);
    }

    const makeRequest = async (requestOptions) => {
        // @ts-ignore
        const res = await fetch(url, options);
        return res;
    }

    // Make the initial request
    let res = await makeRequest(options);

    // Handle 401 Unauthorized by refreshing token
    if (res.status === 401) {
        const tokenRefreshed = await refreshAccessToken();
        if (tokenRefreshed) {
            console.log("Access token refreshed successfully");
            // Retry the request after refreshing token
            res = await makeRequest(options);
            if (res.status === 401) {
                throw new Error(REFRESH_ERROR_MESSAGE);
            }
        } else {
            throw new Error(REFRESH_ERROR_MESSAGE);
        }
    }

    // Check for HTTP errors
    if (!res.ok) {
        throw new Error(await res.text());
    }

    // Try to parse JSON, but let pages handle 204/empty responses as needed
    try {
        return await res.json();
    } catch {
        return null;
    }
}

export async function refreshAccessToken() {
    const url = `${BASE_URL}/refresh`;
    const options = {
        method: "POST",
        credentials: "include" // needed for cookies/auth
    };
    try {
        // @ts-ignore
        const res = await fetch(url, options);
        if (!res.ok) {
            console.error("Failed to refresh access token:", await res.text());
            return false;
        }

        return true;
    } catch (error) {
        console.error("Error refreshing access token:", error);
        return false;
    }
}