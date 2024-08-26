use anyhow::{Context, Result};
use librqbit::{AddTorrent, AddTorrentOptions, AddTorrentResponse, Session};
use reqwest::get;
use std::time::Duration;
use tracing::info;

#[tokio::main]
async fn main() -> Result<()> {
    let output_dir = "data";

    // https://annas-archive.li/md5/08b0f97b98c977da93cd5e5623686af5
    let only_files = "08b0f97b98c977da93cd5e5623686af5";
    let torrent = "https://annas-archive.li/dyn/small_file/torrents/external/libgen_rs_non_fic/r_4319000.torrent";

    // Fetch torrent
    let response = get(torrent).await?.error_for_status()?;
    let content = response.bytes().await?;

    // Output logs to console.
    tracing_subscriber::fmt::try_init().unwrap_or(());

    // Create the session
    let session = Session::new(output_dir.into())
        .await
        .context("error creating session")?;

    // Add the torrent to the session
    let handle = match session
        .add_torrent(
            AddTorrent::from_bytes(content.as_ref()),
            Some(AddTorrentOptions {
                overwrite: true,
                only_files_regex: Some(only_files.to_string()),
                ..Default::default()
            }),
        )
        .await
        .context("error adding torrent")?
    {
        AddTorrentResponse::Added(_, handle) => handle,
        _ => unreachable!(),
    };

    // Print stats periodically.
    tokio::spawn({
        let handle = handle.clone();
        async move {
            let duration = Duration::from_secs(1);
            loop {
                tokio::time::sleep(duration).await;
                let stats = handle.stats();
                info!("{stats}");
            }
        }
    });

    // Wait until the download is completed
    handle.wait_until_completed().await?;

    info!("torrent downloaded");
    Ok(())
}
