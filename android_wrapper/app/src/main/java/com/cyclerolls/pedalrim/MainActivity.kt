package com.cyclerolls.pedalrim

import android.os.Bundle
import android.view.ViewGroup
import androidx.appcompat.app.AppCompatActivity
import com.cyclerolls.pedalrim.mobile.EbitenSurfaceView

class MainActivity : AppCompatActivity() {
    private var ebitenView: EbitenSurfaceView? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Initialize and attach Ebiten SurfaceView to activity layout
        ebitenView = EbitenSurfaceView(this)
        setContentView(
            ebitenView,
            ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT
            )
        )
    }

    override fun onPause() {
        super.onPause()
        ebitenView?.onPause()
    }

    override fun onResume() {
        super.onResume()
        ebitenView?.onResume()
    }
}
