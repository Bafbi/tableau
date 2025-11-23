#!/bin/bash
set -e

# Ensure we are using the dev binary
TABLEAU="./bin/tableau"
# TABLEAU_DIR is set by mise, but we default it here just in case
export TABLEAU_DIR="${TABLEAU_DIR:-.tableau_dev}"

echo "🎨 Setting up demo configuration..."
mkdir -p "$TABLEAU_DIR"
cat > "$TABLEAU_DIR/config.toml" <<EOF
[project]
name = "Tableau Demo"

[columns]
todo = "Backlog"
doing = "In Flight"
review = "Review"
done = "Shipped"

[style]
border_color = "#00FFFF"
selected_color = "#FF00FF"
EOF

echo "📝 Creating tasks..."
$TABLEAU new "Refactor authentication middleware"
$TABLEAU new "Update documentation for API v2"
$TABLEAU new "Fix memory leak in worker pool"
$TABLEAU new "Design new dashboard layout"
$TABLEAU new "Setup Prometheus metrics"
$TABLEAU new "Investigate slow database queries"
$TABLEAU new "Upgrade Go version to 1.22"
$TABLEAU new "Add dark mode support"

echo "🔒 Blocking tasks..."
$TABLEAU block 3

echo "🚀 Simulating active state (modifying files directly)..."
# Move some to Doing (In Flight)
# Note: We use direct file manipulation to avoid creating git branches for this demo
sed -i 's/status = "todo"/status = "doing"/' "$TABLEAU_DIR/tasks/1-refactor-authentication-middleware.md"
sed -i 's/status = "todo"/status = "doing"/' "$TABLEAU_DIR/tasks/4-design-new-dashboard-layout.md"

# Move some to Done (Shipped)
sed -i 's/status = "todo"/status = "done"/' "$TABLEAU_DIR/tasks/2-update-documentation-for-api-v2.md"
sed -i 's/status = "todo"/status = "done"/' "$TABLEAU_DIR/tasks/7-upgrade-go-version-to-1.22.md"

# Set some priorities
sed -i 's/priority = "medium"/priority = "high"/' "$TABLEAU_DIR/tasks/3-fix-memory-leak-in-worker-pool.md"
sed -i 's/priority = "medium"/priority = "low"/' "$TABLEAU_DIR/tasks/8-add-dark-mode-support.md"

echo "✅ Demo state ready!"
$TABLEAU list
