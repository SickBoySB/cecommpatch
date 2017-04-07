run("stats")
print("gameinit executed\n" )


spatialDumpOn = false

function toggleSpatialDump()

	spatialDumpOn = not spatialDumpOn
		
	VALUE_STORE[ "SpatialSVGDump" ] = spatialDumpOn
	VALUE_STORE[ "PathfindingSVGDump" ] = spatialDumpOn

end

VALUE_STORE[ "ScriptDebugEnabled" ] = true

function toggleScriptDebug()
	VALUE_STORE[ "ScriptDebugEnabled" ] = not VALUE_STORE[ "ScriptDebugEnabled" ]
end

function toggleFastSimMode()
	VALUE_STORE[ "2xSimMode" ] = not VALUE_STORE[ "2xSimMode" ]
end

function toggleJPS()
	VALUE_STORE[ "UseJumpPointSearch" ] = not VALUE_STORE[ "UseJumpPointSearch" ]
end

function spawnEvent( name )
	local eventQ = query("gameSimEventManager", "startEvent", name, {}, {})
	return eventQ[ 1 ]
end

registerHotKey( true, true, false, "s", "toggleSpatialDump()" )
registerHotKey( true, true, false, "t", "toggleTelemetry()" )
registerHotKey( true, true, false, "d", "toggleScriptDebug()" )
registerHotKey( true, true, false, "f", "toggleFastSimMode()" )
registerHotKey( true, true, false, "j", "toggleJPS()" )

-- Debugging
VALUE_STORE[ "SkipDescriptionCheck" ] = false
VALUE_STORE[ "DisablePause" ] = false

-- Gameplay 
VALUE_STORE[ "IgnoreBuildingCosts" ] = false
VALUE_STORE[ "FastSimMode" ] = false
VALUE_STORE[ "PopulateFauna" ] = true
VALUE_STORE[ "PopulateHorrors" ] = false

-- Path finding parameters 
VALUE_STORE[ "DefaultPathWaiting" ] = 1
VALUE_STORE[ "AStarOccupancyCost" ] = 50
VALUE_STORE[ "AStarSqueezeCost" ] = 50
VALUE_STORE[ "AStarCheckOccupancyRadius" ] = 9
VALUE_STORE[ "AStarDirectionChangeCost" ] = 0
VALUE_STORE[ "WeakFollowDefault" ] = 5
VALUE_STORE[ "UseJumpPointSearch" ] = false
VALUE_STORE[ "MaxJPSDistance" ] = 25
VALUE_STORE[ "MaximumPathfindingRetries" ] = 4 --50

-- Testing 
VALUE_STORE[ "DeterminismFailureFatal" ] = false
VALUE_STORE[ "NetworkDeterminismTest" ] = true
VALUE_STORE[ "SerializationStressTest" ] = false
VALUE_STORE[ "SerializationStressTestActuallySave" ] = false


-- Debug output 
VALUE_STORE[ "DeterminismStateDump" ] = false
VALUE_STORE[ "DeterminismStateDumpToConsole" ] = false
VALUE_STORE[ "DeterminismStateDumpScriptState" ] = false
VALUE_STORE[ "DeterminismStateDumpInterval" ] = false
VALUE_STORE[ "DeterminismStateDumpMinFrame" ] = 0
VALUE_STORE[ "DeterminismStateDumpMaxFrame" ] = 5000
VALUE_STORE[ "DeterminismStateModuloFrame" ] = 1
VALUE_STORE[ "StoreReplay" ] = true
VALUE_STORE[ "SpatialSVGDumpOnFailure" ] = false
VALUE_STORE[ "SpatialSVGDumpOnConnectedUpdate" ] = false
VALUE_STORE[ "SpatialSVGPerformanceDump" ] = false 
VALUE_STORE[ "pathingPerformanceWarningThreshold" ] = 30 
VALUE_STORE[ "PathfindingSVGDump" ] = false
VALUE_STORE[ "SpatialSVGDump" ] = false

VALUE_STORE[ "SpatialSVGWritePassabilityLandscape" ] = true
VALUE_STORE[ "SpatialSVGWritePassabilityWater" ] = true
VALUE_STORE[ "SpatialSVGWritePassabilityObject" ] = true
VALUE_STORE[ "SpatialSVGWritePassabilityObjectBorder" ] = true
VALUE_STORE[ "SpatialSVGWriteCivilization" ] = false
VALUE_STORE[ "SpatialSVGWriteNudge" ] = true
VALUE_STORE[ "SpatialSVGWriteOccupancy" ] = true
VALUE_STORE[ "SpatialSVGWriteGeneratedPaths" ] = true
VALUE_STORE[ "SpatialSVGWriteAccessPoints" ] = false
VALUE_STORE[ "SpatialSVGWriteDebugMarkers" ] = true

VALUE_STORE[ "StatsAverageWindow" ] = 30
VALUE_STORE[ "TimingCSVLog" ] = false
--VALUE_STORE[ "TimingCSVLogDir" ] = "C:\\Users\\MjB\\Google Drive\\stats\\"

VALUE_STORE[ "VerboseFSM" ] = false
VALUE_STORE[ "ConsoleOldAStarResults" ] = false
VALUE_STORE[ "DisplayCharLocations" ] = false
VALUE_STORE[ "DisplayDebugOverlay" ] = true
VALUE_STORE[ "ShowFPS" ] = false
VALUE_STORE[ "ExitOnTestComplete" ] = false
VALUE_STORE[ "TestMaxFrames" ] = 500

VALUE_STORE["showCombatDebugConsole"] = false
VALUE_STORE["showCitizenDebugConsole"] = false
VALUE_STORE["showBuildingDebugConsole"] = false
VALUE_STORE["showModuleDebugConsole"] = false
VALUE_STORE["showAnimalDebugConsole"] = false
VALUE_STORE["showFarmDebugConsole"] = false
VALUE_STORE["showFSMDebugConsole"] = true




