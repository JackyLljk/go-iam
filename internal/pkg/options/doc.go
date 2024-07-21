/*
internal/pkg/options 是项目内通用 api 服务器使用的公共标志和选项。
它需要一组最小的依赖项，并且不引用实现，以确保它可以被多个组件（例如希望生成或验证配置的 CLI 命令）重用。
*/

package options
