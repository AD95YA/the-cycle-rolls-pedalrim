package com.cyclerolls.pedalrim

import android.os.Bundle
import android.view.ViewGroup
import androidx.appcompat.app.AppCompatActivity
import com.cyclerolls.pedalrim.mobile.EbitenView

class MainActivity : AppCompatActivity() {
    private var ebitenView: EbitenView? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Initialize and attach public EbitenView to activity layout
        val view = EbitenView(this)
        ebitenView = view
        setContentView(
            view,
            ViewGroup.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT,
                ViewGroup.LayoutParams.MATCH_PARENT
            )
        )
    }

    override fun onPause() {
        super.onPause()
        ebitenView?.suspendGame()
    }

    override fun onResume() {
        super.onResume()
        ebitenView?.resumeGame()
    }
}
