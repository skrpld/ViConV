package com.example.client

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform