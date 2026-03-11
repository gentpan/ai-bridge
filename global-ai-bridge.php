<?php
/**
 * Plugin Name: AI Bridge
 * Plugin URI: https://xifeng.net/global-ai-bridge-wp-plugins
 * Description: Connect WordPress AI features to global AI providers through your own proxy gateway.
 * Version: 1.1.0
 * Author: 西风
 * Author URI: https://xifeng.net
 * Text Domain: global-ai-bridge
 * Requires at least: 6.0
 * Requires PHP: 7.4
 *
 * @package GlobalAIBridge
 */

if ( ! defined( 'ABSPATH' ) ) {
	exit;
}

define( 'GAB_VERSION', '1.1.0' );
define( 'GAB_PLUGIN_FILE', __FILE__ );
define( 'GAB_PLUGIN_DIR', plugin_dir_path( __FILE__ ) );
define( 'GAB_PLUGIN_URL', plugin_dir_url( __FILE__ ) );

require_once GAB_PLUGIN_DIR . 'includes/class-gab-logger.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-proxy-client.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-http-interceptor.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-connectors-compat.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-theme-compat.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-settings.php';
require_once GAB_PLUGIN_DIR . 'includes/class-gab-plugin.php';

register_activation_hook( __FILE__, array( 'GAB_Settings', 'initialize_defaults' ) );

/**
 * Returns the plugin singleton.
 *
 * @return GAB_Plugin
 */
function gab_ai_bridge() {
	return GAB_Plugin::instance();
}

gab_ai_bridge();

/**
 * Convenience wrapper for sending AI requests through the configured proxy.
 *
 * @param array $messages Chat-style message payload.
 * @param array $args     Additional request arguments.
 * @return array|WP_Error
 */
function gab_send_ai_request( array $messages, array $args = array() ) {
	return gab_ai_bridge()->get_proxy_client()->send_chat_request( $messages, $args );
}
