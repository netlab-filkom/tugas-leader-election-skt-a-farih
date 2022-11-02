# Tugas Leader Election Raft

## Kelompok
1. 205150301111026 Farih Akmal Haqiqi
2. 205150300111047 Muhammad Daffa Pradipta Akbar

## Deskripsi Tugas
Implementasikan mekanisme pemilihan leader di protokol Raft dengan alur sebagai berikut :

1. Semua server/node/peer mengawali state sebagai Follower.
2. Sebagai Follower, setiap server memiliki countdown timer. Timer akan berjalan sampai menerima heartbeat message RPC AppendEntries dari leader. Countdown timer akan diset ulang dengan nilai random ketika pertama kali dijalankan atau ketika menerima heartbeat dari Leader.
3. Server yang countdown timernya habis akan berubah state menjadi Candidate. Sebagai Candidate, server melakukan hal beriut :
    - Menaikkan term-nya 1 poin
    - Memberikan 1 vote ke dirinya sendiri
    - Menyebarkan RPC RequestVote untuk meminta voting dari semua server lain. Server lain wajib membalas RequestVote jika term Candidate lebih tinggi dari term Follower.
    - Candidate berubah menjadi Leader jika mendapat suara > 50% jumlah server.

## Ketentuan Tugas
1. Setiap server menyimpan informasi :
    - ID, alamat IP dan port dari semua server.
    - Term dirinya sendiri
2. Countdown timer di-set secara acak dengan periode 30-50 detik
3. Jeda pengiriman heartbeat message AppendEntries 10 detik
4. Komunikasi antar server menggunakan mekanisme RPC dengan dua remote function :
    - Fungsi AppendEntries(term, index, command). Untuk sementara index = 0, command = "", dan term disesuaikan dengan term Leader.
    - Fungsi RequestVote() -> True/False. Fungsi request vote ke server lain. Return value True jika server melakukan vote dan False jika server tidak melakukan vote. 
