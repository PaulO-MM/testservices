// Standard error response format shared across the Node.js API.
module.exports = {
  successResponse(data) {
    return { success: true, data };
  },

  errorResponse(code, message, details = undefined) {
    const resp = {
      success: false,
      error: { code, message },
    };
    if (details !== undefined) {
      resp.error.details = details;
    }
    return resp;
  },
};
