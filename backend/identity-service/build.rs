fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_prost_build::configure().compile_protos(
        &[
            "proto/authentication.proto",
            "proto/user.proto",
            "proto/authorization/authorization.proto",
            "proto/authorization/permission.proto",
            "proto/authorization/role.proto",
            "proto/authorization/user_role.proto",
        ],
        &["proto"],
    )?;
    Ok(())
}
