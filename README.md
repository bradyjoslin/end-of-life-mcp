# 📦 End-of-Life Info Tool

This Go [mcp.run](https://www.mcp.run/bradyjoslin/end-of-life) servlet provides a simple API-style interface for interacting with [endoflife.date](https://endoflife.date) data. It allows you to query available products, their release cycles, and detailed end-of-life (EOL) information for specific product versions.

## 🔧 Available Tools

### 1. `list_available_products`
Lists all software, operating systems, and devices tracked by endoflife.date.

#### Description
Use this to discover what products you can query. Typical examples include `"ubuntu"`, `"php"`, `"windows"`, etc.

#### Parameters
_None_

### 2. `get_product_cycles`
Returns a list of release cycles for a given product, along with their end-of-life (EOL) dates.

#### Parameters

```sh
{
  "product_name": "string" // e.g. "ubuntu", "nodejs"
}
```

### 3. `get_cycle_details`
Returns detailed information about a specific product release cycle, such as support dates, LTS status, and latest version.

#### Parameters

```sh
{
  "product_name": "string", // e.g. "ubuntu"
  "cycle_name": "string"    // e.g. "22.04"
}
```
