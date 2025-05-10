package driver

// QueryCondition represents a collection of query conditions.
// Contains optional time range, node conditions, and multiple value range filters.
type QueryCondition struct {
	*ConTimeRange               // Embedded time range condition (optional)
	*ConNode                    // Embedded node condition (optional)
	ValRange      []ConValRange // Slice of value range conditions
}

// ConTimeRange defines time range query conditions.
// Uses pointers to implement optional fields (nil means no time constraint).
type ConTimeRange struct {
	StartAt *uint64 // Start timestamp in UNIX format (inclusive, closed interval)
	EndAt   *uint64 // End timestamp in UNIX format (inclusive, closed interval)
}

// ConNode defines node query conditions.
// Pointer fields allow empty values, both nil means no node constraint.
type ConNode struct {
	NodeName *string // Node name for exact match
	NodeID   *string // Node ID for exact match
}

// ConValRange defines value range filter conditions.
// Used to specify range constraints for specific keys.
// Support int and float data.
type ConValRange struct {
	Key      string   // Target key name for filtering
	StartVal *float64 // Minimum value pointer (inclusive, nil means no lower bound)
	EndVal   *float64 // Maximum value pointer (inclusive, nil means no upper bound)
}

// ConditionOption defines function type for functional option pattern
// Used to flexibly compose query conditions
type ConditionOption func(driver *QueryCondition)

// WithTimeRange creates a time range condition configuration
// Parameters:
//   - startAt: UNIX timestamp in seconds (0 means not set)
//   - endAt: UNIX timestamp in seconds (0 means not set)
//
// Creates time condition only when startAt/endAt are non-zero
func WithTimeRange(startAt, endAt uint64) ConditionOption {
	return func(driver *QueryCondition) {
		if startAt != 0 {
			driver.StartAt = &startAt
		}
		if endAt != 0 {
			driver.EndAt = &endAt
		}
	}
}

// WithNode creates a node condition configuration
// Parameters:
//   - nodeName: Exact match node name (empty string means not set)
//   - nodeID: Exact match node ID (empty string means not set)
//
// Creates node condition only when either parameter is non-empty
func WithNode(nodeName, nodeID string) ConditionOption {
	return func(driver *QueryCondition) {
		if nodeName != "" {
			driver.NodeName = &nodeName
		}
		if nodeID != "" {
			driver.NodeID = &nodeID
		}
	}
}

// WithValRange creates a value range condition configuration
// Parameters:
//   - key: Required key name for filtering
//   - startVal: Minimum value (inclusive, 0 is considered valid)
//   - endVal: Maximum value (inclusive, 0 is considered valid)
//
// Appends new condition to existing ValRange slice
func WithValRange(key string, startVal, endVal float64) ConditionOption {
	return func(driver *QueryCondition) {
		driver.ValRange = append(driver.ValRange, ConValRange{
			Key:      key,
			StartVal: &startVal,
			EndVal:   &endVal,
		})
	}
}
