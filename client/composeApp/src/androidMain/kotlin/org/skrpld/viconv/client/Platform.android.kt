package org.skrpld.viconv.client

import android.os.Build
import com.skrpld.nearbeee.client.Platform

class AndroidPlatform : Platform {
    override val name: String = "Android ${Build.VERSION.SDK_INT}"
}

actual fun getPlatform(): Platform = AndroidPlatform()