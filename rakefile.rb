task :default => [:help]

def encrypted_file
    return File.expand_path "~/Documents/t/t.zip"
end

def decrypted_folder_parent
    return "~/decrypted_folder/"
end

task :help do
    puts "TreuCrypt replacement"
    puts "encrypted_file          => #{encrypted_file}"
    puts "decrypted_folder_parent => #{decrypted_folder_parent}"
end
