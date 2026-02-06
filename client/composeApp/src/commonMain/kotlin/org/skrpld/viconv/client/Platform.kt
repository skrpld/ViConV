package org.skrpld.viconv.client

interface Platform {
    val name: String
}

expect fun getPlatform(): Platform