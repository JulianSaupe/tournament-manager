fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_prost_build::configure().compile_protos(
        &[
            "../shared/proto/authentication.proto",
            "../shared/proto/user.proto",
            "../shared/proto/authorization/authorization.proto",
            "../shared/proto/authorization/permission.proto",
            "../shared/proto/authorization/role.proto",
            "../shared/proto/authorization/user_role.proto",
        ],
        &["../shared/proto"],
    )?;
    Ok(())
}
