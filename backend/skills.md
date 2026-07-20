# AI Agent Skills Specification

[
  {
    "name": "get_stock_price",
    "description": "ใช้สำหรับดึงราคาหุ้นล่าสุด โดยระบุชื่อย่อหุ้น เช่น AAPL, GOOG",
    "input_schema": {
      "type": "object",
      "properties": {
        "ticker": {
          "type": "string",
          "description": "ชื่อย่อของหุ้น (Ticker symbol) เช่น AAPL"
        }
      },
      "required": ["ticker"]
    }
  }
]