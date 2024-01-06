const fs = require('fs')
module.exports = ({github, context}) => {
    const data = fs.readFileSync('output.sh', 'utf8')
    const body = "Hello，您可以直接执行该命令：\n" +
        "\n" +
        "```shell\n" +
        data +
        "\n" +
        "```\n" +
        "\n" +
        "希望可以帮助到您"
    return body
}