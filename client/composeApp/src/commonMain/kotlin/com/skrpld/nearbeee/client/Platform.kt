package com.skrpld.nearbeee.client

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform