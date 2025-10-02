// Client-side logging utility with structured logging
class Logger {
  constructor() {
    this.isProduction = process.env.NODE_ENV === 'production';
    this.logLevel = process.env.NEXT_PUBLIC_LOG_LEVEL || (this.isProduction ? 'error' : 'debug');
  }

  // Log levels: debug < info < warn < error
  levels = {
    debug: 0,
    info: 1,
    warn: 2,
    error: 3,
  };

  shouldLog(level) {
    return this.levels[level] >= this.levels[this.logLevel];
  }

  formatMessage(level, message, context = {}) {
    const timestamp = new Date().toISOString();
    const logEntry = {
      timestamp,
      level: level.toUpperCase(),
      message,
      ...context,
    };

    // Add user context if available
    if (typeof window !== 'undefined' && window.user) {
      logEntry.user_id = window.user.sub;
    }

    // Add session context
    if (typeof window !== 'undefined') {
      logEntry.session_id = this.getSessionId();
      logEntry.url = window.location.href;
      logEntry.user_agent = navigator.userAgent;
    }

    return logEntry;
  }

  getSessionId() {
    if (typeof window === 'undefined') return null;
    
    let sessionId = sessionStorage.getItem('session_id');
    if (!sessionId) {
      sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
      sessionStorage.setItem('session_id', sessionId);
    }
    return sessionId;
  }

  debug(message, context = {}) {
    if (!this.shouldLog('debug')) return;
    
    const logEntry = this.formatMessage('debug', message, context);
    console.debug('[DEBUG]', logEntry);
  }

  info(message, context = {}) {
    if (!this.shouldLog('info')) return;
    
    const logEntry = this.formatMessage('info', message, context);
    console.info('[INFO]', logEntry);
  }

  warn(message, context = {}) {
    if (!this.shouldLog('warn')) return;
    
    const logEntry = this.formatMessage('warn', message, context);
    console.warn('[WARN]', logEntry);
    
    // Send to monitoring service in production
    if (this.isProduction) {
      this.sendToMonitoring('warn', logEntry);
    }
  }

  error(message, context = {}) {
    if (!this.shouldLog('error')) return;
    
    const logEntry = this.formatMessage('error', message, context);
    console.error('[ERROR]', logEntry);
    
    // Always send errors to monitoring
    this.sendToMonitoring('error', logEntry);
  }

  // Send logs to external monitoring service
  sendToMonitoring(level, logEntry) {
    if (typeof window === 'undefined') return;

    // Example: Send to monitoring service
    // Replace with your actual monitoring service (Sentry, LogRocket, etc.)
    try {
      // This would typically be an async call to your monitoring service
      if (window.reportError && level === 'error') {
        window.reportError(new Error(logEntry.message), logEntry);
      }
      
      // Or send to custom endpoint
      // fetch('/api/logs', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(logEntry),
      // }).catch(() => {}); // Fail silently for logging
    } catch (error) {
      // Fail silently - don't let logging errors break the app
      console.error('Failed to send log to monitoring:', error);
    }
  }

  // Performance timing helper
  time(label) {
    if (!this.shouldLog('debug')) return;
    console.time(label);
  }

  timeEnd(label, context = {}) {
    if (!this.shouldLog('debug')) return;
    console.timeEnd(label);
    this.debug(`Timer ${label} completed`, context);
  }

  // API call logging helper
  logApiCall(method, url, status, duration, context = {}) {
    const message = `API ${method} ${url} - ${status} (${duration}ms)`;
    const logContext = {
      api_method: method,
      api_url: url,
      api_status: status,
      api_duration: duration,
      ...context,
    };

    if (status >= 400) {
      this.error(message, logContext);
    } else if (status >= 300) {
      this.warn(message, logContext);
    } else {
      this.info(message, logContext);
    }
  }
}

// Create singleton instance
export const logger = new Logger();

// Export for default import
export default logger;