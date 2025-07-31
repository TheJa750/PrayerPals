// src/lib/api.js
export const BASE_URL = "http://localhost:8080/api"; // Adjust when deploying

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
    // @ts-ignore
    const res = await fetch(url, options);
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