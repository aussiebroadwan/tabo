const LOG_LEVELS = ['DEBUG', 'INFO', 'WARN', 'ERROR'];
const CURRENT_LOG_LEVEL = import.meta.env.VITE_LOG_LEVEL || 'INFO';
const LOG_LEVEL_INDEX = LOG_LEVELS.indexOf(CURRENT_LOG_LEVEL);

export default function log(level, message, context = {}) {
    if (LOG_LEVELS.indexOf(level) >= LOG_LEVEL_INDEX) {
        console.log(JSON.stringify({
            timestamp: new Date().toISOString(),
            level,
            message,
            context
        }));
    }
}
