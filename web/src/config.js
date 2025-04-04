// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Define base API paths based on environment
const BASE_API = import.meta.env.PROD
    ? "/api"
    : "http://localhost:8181/api";

export const API_URL = `${BASE_API}/books`;
export const LOGIN_URL = `${BASE_API}/login`;
