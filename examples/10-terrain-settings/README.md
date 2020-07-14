# Terrain generator

This application generates a terrain. The TerrainBuilder is managable from the settings menu.

## Settings menu

- **Rows (i)** - It is the input of the TerrainBuilder.SetWidth function.
- **Cols (i)** - It is the input of the TerrainBuilder.SetLength function.
- **Iter (i)** - It is the input of the TerrainBuilder.SetIterations function.
- **Scale (f)** - It is the input of the TerrainBuilder.SetScale function. (X, Z axis)
- **Peak (i)** - It is the input of the TerrainBuilder.SetPeakProbability function.
- **Cliff (i)** - It is the input of the TerrainBuilder.SetCliffProbability function.
- **MinH (f)** - It is the input of the TerrainBuilder.SetMinHeight function.
- **MaxH (f)** - It is the input of the TerrainBuilder.SetMaxHeight function.
- **PosY (f)** - It is the input of the TerrainBuilder.SetPosition function. (Y axis)
- **Seed (i)** - It is the input of the TerrainBuilder.SetSeed function, if the **RandSeed** flag is not set.
- **RandSeed** - If this flag is set, it calls the TerrainBuilder.RandomSeed function instead of SetSeed.
- **Terr tex** - It is the (string) enum of the surface texture. Currently only the 'Grass' is supported.
- **HasLiquid** - If this flag is set, the terrain is generated with liquid.
- **Leta (f)** - It is the input of the TerrainBuilder.SetLiquidEta function.
- **Lampl (f)** - It is the input of the TerrainBuilder.SetLiquidAmplitude function.
- **Lfreq (f)** - It is the input of the TerrainBuilder.SetLiquidFrequency function.
- **Ldetail (i)** - It is the input of the TerrainBuilder.SetLiquidDetailMultiplier function.
- **W lev (f)** - It is the input of the TerrainBuilder.SetLiquidWaterLevel function.
- **Liq tex** It is the (string) enum of the surface liquid. Currently only the 'Water' is supported.
- **Debug model** - It turns the debug mode on or off.
- **Bg R (f)** - It is the red component of the background color [0-1]
- **Bg G (f)** - It is the green component of the background color [0-1]
- **Bg B (f)** - It is the blue component of the background color [0-1]
